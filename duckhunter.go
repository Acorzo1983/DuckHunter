package main

import (
    "bufio"
    "flag"
    "fmt"
    "net"
    "os"
    "strings"
    "sync"
)

const version = "0.9"

// ANSI escape codes for colors
const (
    Red    = "\033[31m"
    Green  = "\033[32m"
    Orange = "\033[33m"
    Reset  = "\033[0m"
)

// Vulnerable patterns for subdomain takeover
var vulnerablePatterns = map[string][]string{
    "No A records found": {
        "digitalocean.com",
        "aws.amazon.com",
        "github.io",
        "herokuapp.com",
        "pantheon.io",
        "bitbucket.io",
        "fastly.net",
        "ghost.io",
        "wordpress.com",
    },
    "NXDOMAIN": {
        "unconfigured.herokudns.com",
        "namecheap.com",
        "myshopify.com",
        "cloudapp.net",
        "smugmug.com",
        "cargo.site",
        "pageserve.co",
        "domains.goog",
        "azurewebsites.net",
    },
}

// Resolve DNS for the given subdomain
func resolveDNS(subdomain string, wg *sync.WaitGroup, results chan<- string) {
    defer wg.Done()
    fmt.Printf("Resolving DNS for: %s\n", subdomain)
    ips, err := net.LookupIP(subdomain)
    if err != nil {
        fmt.Printf(Red+"Failed to resolve: %s\n"+Reset, subdomain)
        results <- fmt.Sprintf("%s,No A records found or error", subdomain)
        return
    }
    fmt.Printf("Resolved: %s to IPs: %v\n", subdomain, ips)
    results <- fmt.Sprintf("%s,%v", subdomain, ips)
}

// Check if the subdomain is vulnerable based on DNS resolution response
func checkVulnerability(subdomain string, response string) (bool, string) {
    for pattern, providers := range vulnerablePatterns {
        if strings.Contains(response, pattern) {
            for _, provider := range providers {
                if strings.Contains(subdomain, provider) {
                    fmt.Printf(Green+"Vulnerability found for: %s\n"+Reset, subdomain)
                    return true, fmt.Sprintf("%s is vulnerable under %s (%s)", subdomain, provider, pattern)
                }
            }
        }
    }
    fmt.Printf(Orange+"No vulnerability found for: %s\n"+Reset, subdomain)
    return false, fmt.Sprintf("[Not Vulnerable] %s", subdomain)
}

// Generate the report and save it to a file
func generateReport(results []string, outputFile string) {
    fmt.Printf("Generating report: %s\n", outputFile)
    file, err := os.Create(outputFile)
    if err != nil {
        fmt.Println("Error creating the report file:", err)
        return
    }
    defer file.Close()

    for _, result := range results {
        file.WriteString(result + "\n")
    }

    fmt.Printf("Analysis completed. Check the %s file for results.\n", outputFile)
}

func main() {
    // Custom Usage function
    flag.Usage = func() {
        fmt.Fprintf(flag.CommandLine.Output(), "Usage of Duckhunter:\n")
        fmt.Fprintf(flag.CommandLine.Output(), "  -l <file>    : Specify a list of subdomains to check\n")
        fmt.Fprintf(flag.CommandLine.Output(), "  -d <domain>  : Specify a single domain to check\n")
        fmt.Fprintf(flag.CommandLine.Output(), "  -o <output>  : Specify the output file (default: report.txt)\n")
        fmt.Fprintf(flag.CommandLine.Output(), "  -h           : Display this help message\n")
    }

    // Print version information and credits in orange
    fmt.Printf("DuckHunter version %s - " + Orange + "Made with love by Albert.C" + Reset + "\n\n", version)

    // Define and parse flags
    subdomainList := flag.String("l", "", "List of subdomains to check")
    singleDomain := flag.String("d", "", "Single domain to check")
    outputFile := flag.String("o", "report.txt", "Output file")
    flag.Parse()

    if len(os.Args) == 1 {
        flag.Usage()
        return
    }

    var subdomains []string

    // Process subdomains list
    if *subdomainList != "" {
        file, err := os.Open(*subdomainList)
        if err != nil {
            fmt.Println("Error reading file:", err)
            return
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            subdomains = append(subdomains, scanner.Text())
        }
    }

    // Process single domain
    if *singleDomain != "" {
        subdomains = append(subdomains, *singleDomain)
    }

    if len(subdomains) == 0 {
        fmt.Println("Error: No subdomains or domain provided.")
        return
    }

    var wg sync.WaitGroup
    results := make(chan string)

    go func() {
        wg.Wait()
        close(results)
    }()

    // Resolve DNS and check vulnerability
    for _, subdomain := range subdomains {
        wg.Add(1)
        go resolveDNS(subdomain, &wg, results)
    }

    var reportResults []string
    for result := range results {
        subdomain := result[:strings.Index(result, ",")]
        response := result[strings.Index(result, ",")+1:]
        if strings.Contains(response, "No A records found or error") {
            reportResults = append(reportResults, fmt.Sprintf("[Not Reachable] %s", subdomain))
        } else {
            vulnerable, message := checkVulnerability(subdomain, response)
            if vulnerable {
                reportResults = append(reportResults, fmt.Sprintf("[Vulnerable] %s", message))
            } else {
                reportResults = append(reportResults, fmt.Sprintf("[Not Vulnerable] %s", message))
            }
        }
    }

    // Generate report
    generateReport(reportResults, *outputFile)
}

# DuckHunter

DuckHunter is a robust tool specifically designed to detect subdomain takeover vulnerabilities, particularly those related to the "Sitting Duck" vulnerability that has put over 1 million domains at risk. By resolving DNS for specified subdomains and analyzing them for patterns associated with potential takeovers, DuckHunter generates detailed reports that help you secure your domains. Whether you're checking a single domain or an extensive list, DuckHunter's color-coded output ensures that results are easy to understand and act upon.

## Features

- **Version Information**: Displays the current version and attribution at the start.
- **DNS Resolution**: Resolves DNS for each subdomain and checks for A records.
- **Vulnerability Detection**: Identifies vulnerable subdomains based on known patterns.
- **Color-Coded Output**:
  - **Red** for DNS resolution failures (`[Not Reachable]`).
  - **Green** for vulnerable subdomains (`[Vulnerable]`).
  - **Orange** for subdomains that are not vulnerable (`[Not Vulnerable]`).
- **Custom Help Message**: Displays "Usage of Duckhunter:" instead of the default Go usage message.

## Installation

Clone the repository and navigate to the project directory:

 ```bash
git clone https://github.com/Acorzo1983/DuckHunter.git
cd DuckHunter
 ```

## Usage

You can run DuckHunter with various options:

### Help

To display the help message:

```bash
go run duckhunter.go -h
```

### Checking a List of Subdomains

To check a list of subdomains provided in a file:

```bash
go run duckhunter.go -l subdomains.txt -o report.txt
```

### Checking a Single Domain

To check a single domain:

```bash
go run duckhunter.go -d example.com -o report.txt
```

### Example: Using DuckHunter with Subfinder

Subfinder is a subdomain discovery tool that can be used to generate a list of subdomains. Here’s how you can use Subfinder with DuckHunter:

1. **Install Subfinder** (if you haven't already):

    ```bash
    go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest
    ```

2. **Find Subdomains** for a given domain (e.g., `example.com`):

    ```bash
    subfinder -d example.com -o subdomains.txt
    ```

3. **Run DuckHunter** to check for vulnerabilities:

    ```bash
    go run duckhunter.go -l subdomains.txt -o report.txt
    ```

## Output

The results are color-coded in the terminal and saved to the specified output file:

- **[Not Reachable]**: Subdomains that could not be resolved.
- **[Vulnerable]**: Subdomains that are identified as vulnerable.
- **[Not Vulnerable]**: Subdomains that are not vulnerable.

## Example

```bash
go run duckhunter.go -l subdomains.txt -o report.txt
```

This command will read the subdomains from \`subdomains.txt\`, check for vulnerabilities, and save the report to \`report.txt\`.

## Contributing

Feel free to submit issues or pull requests if you have suggestions for improvements or new features.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

**DuckHunter version 0.9 - Made with love by Albert.C**

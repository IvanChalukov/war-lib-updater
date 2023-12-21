# war-lib-updater

war-lib-updater is a bash script designed to update Java libraries within a WAR (Web Application Archive) file, leveraging Maven Central for library downloads. It's particularly useful for automating the process of keeping Java dependencies up-to-date in a WAR file.

## Features

- Downloads and updates libraries in a WAR file from Maven Central.
- Handles both local and remote WAR file sources.
- Automatically cleans up the working directory to avoid conflicts.
- Supports multiple library updates in a single run.

## Prerequisites

- Bash
- Maven
- Java JDK 8 or later
- `curl` or `wget` for downloading files
- `unzip` and `jar` tools for handling WAR files

## Usage

To use war-lib-updater, you need to specify the WAR file (either a local path or a URL) and one or more libraries in the `groupId:artifactId:version` format.

```bash
./war-lib-updater.sh <war-file-or-url> <groupId:artifactId:version> [<groupId:artifactId:version> ...]
```
Example:
```bash
./war-lib-updater.sh ./myapp.war org.apache.commons:commons-lang3:3.12 org.slf4j:slf4j-api:1.7.30
```


### Steps Performed by the Script
1. WAR File Handling: Downloads or copies the specified WAR file to the working directory.
2. Library Management: For each specified library, the script:
- Downloads the library from Maven Central.
- Replaces the old version of the library in the WEB-INF/lib directory of the WAR file with the new version.
3. Repackaging: Repackages the updated WAR file.
4. Cleanup: Performs cleanup operations to maintain a clean working environment.

## Important Notes
- The script cleans up the working directory at the start of its execution. Ensure to run it in an isolated directory to avoid accidental data loss.
- It's recommended to backup your original WAR file before using this script.

## Contributing
Contributions to war-lib-updater are welcome! Feel free to submit pull requests or open issues to propose features or report bugs.
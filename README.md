# xctester

A command line tool to extract test summaries & screenshots from Xcode 11 XCResult files.

## Installation

Download latest version from github releases: 

`wget https://github.com/eroshenkoam/xctester/releases/latest/download/xctester`

And make it executable: 

`chmod +x xctester`

## Usage

`xctester <command> <options>`

Below are a few examples of common commands. For further assistance, use the --help option on any command

### Export to Allure2 results

`xctester export /path/to/Test.xcresult`

After that you can generate Allure report by following command: 

`allure serve /path/to/outputDirectory`

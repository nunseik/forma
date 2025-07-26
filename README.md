# FORMA Project Initializer & Scaffolding
<div align="center">
<img src="https://placehold.co/200x200/a855f7/ffffff?text=FORMA" alt="forma logo" width="200"/>
</div>

<p align="center">
<strong>FORMA is a smart CLI tool that automates project scaffolding using powerful, shareable templates.</strong>
<br />
<!-- TODO: Replace YOUR_USERNAME with your actual GitHub username -->
<a href="https://github.com/nunseik/forma/releases"><strong>» Download a Release «</strong></a>
</p>

-----

## Features

  * **Interactive UI**: Simply run `forma new` to launch a friendly terminal UI that guides you through selecting a template and naming your project.
  * **Template Management**: Add new templates directly from Git repositories (`forma add https://..repo.git`), list templates available (`forma list`) or remove ones you no longer need (`forma remove`).
  * **Powerful Templating**: Uses Go's templating engine to inject variables like project name, author, and timestamps into your files.
  * **Automated Hooks**: Each template can define hooks to automatically run commands like `git init` or `npm install` after project creation.
  * **Cross-Platform**: Built with Go to run natively on Windows, macOS, and Linux.

-----

## Installation

### With Go

If you have the Go toolchain installed, you can easily install FORMA with `go install`:

```bash
go install github.com/nunseik/forma@latest
```

### From GitHub Releases

If you don't have Go, you can download a compiled binary for your operating system from the [Releases page](https://github.com/nunseik/forma/releases). After downloading, unzip the file and place the `forma` (or `forma.exe`) binary in an executable path.

-----

## Usage

FORMA works with simple, intuitive commands.

### Create a New Project

To launch the interactive mode, run:

```bash
forma new
```

You can also provide arguments directly:

```bash
forma new <template-name> <project-name> --author "Your Name"
```

### List Available Templates

Shows all templates currently installed.

```bash
forma list
```

### Add a New Template

Add a new template from a Git repository.

```bash
forma add <git-repo-url>
```

**Example:**

```bash
forma add https://github.com/project-starters/go-cli-template.git
```

### Remove a Template

Delete a template from your local machine.

```bash
forma remove <template-name>
```

-----

## Templates

### How to Create a Template

A FORMA template is a simple directory that follows a few rules.

1.  **Create a Directory**: Make a new directory with a name like `my-new-template`.

2.  **Add Template Files**: Inside this directory, place the files and folders that will be the skeleton of your project. You can use Go template syntax like `{{ .ProjectName }}` and `{{ .Author }}`.

    **Example: `main.go`**

    ```go
    // Author: {{ .Author }}
    // Project: {{ .ProjectName }}

    package main

    import "fmt"

    func main() {
        fmt.Println("Hello, from {{ .ProjectName }}!")
    }
    ```

3.  **Create a `template.yaml`**: In the root of the directory, create a `template.yaml` file that describes the template.

    ```yaml
    # template.yaml
    name: "My Go App"
    description: "A simple Go application template."
    hooks:
      post_create:
        - "git init"
        - "go mod init github.com/{{ .Author }}/{{ .ProjectName }}"
        - "go mod tidy"
        - "echo '✅ Project {{ .ProjectName }} initialized.'"
    ```

### `template.yaml` Fields

  * **`name`**: A human-readable name that will be displayed by `forma list`.
  * **`description`**: A short sentence explaining the template's purpose.
  * **`hooks.post_create`**: A list of commands to be executed after the files have been generated. These commands are run in the root directory of the new project.

### Template Storage Location
FORMA stores all user-added templates in a local configuration directory. This allows you to manually add, edit, or back up your templates.

#### Linux
```
cd ~/.config/forma/templates
```
#### macOS
```
cd ~/Library/Application Support/forma/templates
```
#### Windows
```
cd C:\Users\<YourUser>\AppData\Roaming\forma\templates
```

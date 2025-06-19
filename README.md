# my_vcs

# 🗃️ My-VCS: A Lightweight Version Control System

`my-vcs` is a lightweight, version control system written in Go. Inspired by Git, it implements core VCS functionality including repository initialization, committing, branching, status checks, and more.

---

## 📁 Project Structure



my-vcs/ <br>
├── cmd/<br>
│   ├── messages.go             <br>
│   └── my-cli/<br>
│       └── main.go             <br>
├── pkg/<br>
│   ├── git/                    br>
│   │   ├── add.go<br>
│   │   ├── branch.go<br>
│   │   └── ...<br>
│   ├── compress.go             <br>
│   ├── head.go                 <br>
│   └── ...                     <br>
├── go.mod<br>
├── go.sum<br>
<br>

---

## 🚀 Getting Started

### 🔧 Prerequisites

- Go 1.18+ installed

### 📦 Installation

Clone this repository:

```bash
git clone https://github.com/Havok17122004/my_vcs
cd my-vcs/cmd/my-cli
````

Build the CLI:

```bash
go build -o vcs
```

Or run directly:

```bash
go run main.go <command>
```

--- 


## 📚 Commands and Usage

| Command    | Syntax / Example                                                                               | Description                                                                |
| ---------- | ---------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------- |
| `init`     | `./vcs init` <br> `./vcs init my-project`                                                      | Initializes a new VCS repository in the current directory or specified one |
| `add`      | `./vcs add file1.txt file2.txt`                                                                | Adds files to the staging area and compresses them as blobs                |
| `commit`   | `./vcs commit` <br> Prompt:`Enter the commit message:`                                 | Prompts for a commit message in terminal and creates a commit if staged changes exist |
| `status`   | `./vcs status`                                                                                 | Displays changes staged, modified, or untracked                            |
| `branch`   | `./vcs branch` <br> `./vcs branch new-feature`                                                 | Lists all branches <br> Creates new branches                                     |
| `checkout` | `./vcs checkout branch-name` <br> `./vcs checkout hash file.txt`                               | Switches to branch or specific commit (with optional path filtering)       |
| `config`   | `./vcs config user.name "Atharv"` <br> `./vcs config user.email`                                | Set config values like username or email <br> Get config values like username or email                           |
| `diff`     | `./vcs diff` <br> `./vcs diff file.go` <br> `./vcs diff c1 c2` <br> `./vcs diff c1 c2 file.go` | Compares changes between working directory, staging area, or commits       |
| `cat-file` | `./vcs cat-file -p hash` <br> `./vcs cat-file -t hash` <br> `./vcs cat-file -s hash`           | Decompresses objects and shows contents, type, or size                     |
| `log`      | `./vcs log` <br> `./vcs log main ^feature`                                                     | Shows commit history, with support for filters and exclusions              |
| `reset`    | `./vcs reset --soft hash` <br> `--mixed` (default) <br> `--hard`                               | Resets HEAD and optionally staging area or working directory               |
| `merge`    | `./vcs merge branch-name`                                                                      | Merges another branch into the current one                                 |

---

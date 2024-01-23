# create new project

quick create new project from template.

## Flow

1. clone template project to tmp folder
2. copy template project to new project
   - remove `.git` folder in new project
3. rename project name in new project
4. rename some dir name in new project
5. rename some file name in new project
6. render some file in new project
7. remove tmp folder

## Usage

```bash
kite new -c config.yml <project-name>
kite new --from /path/to/tpl-repo -v env=pre <project-name>
```

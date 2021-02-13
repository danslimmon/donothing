# donothing

`donothing` is a Go framework for [do-nothing
scripting](https://blog.danslimmon.com/2019/07/15/do-nothing-scripting-the-key-to-gradual-automation/).
Do-nothing scripting is an approach to task automation that lets your processes evolve with minimal
effort from documented manual processes to fully automated.

A `donothing` script walks through a **procedure**, which comprises a sequence of **steps**. As an
example, here's a simple procedure for restoring a database backup:

```go
package main

import (
  "github.com/danslimmon/donothing"
)

func main() {
  pcd := donothing.NewProcedure()
  pcd.Short(`Restore a database backup`)

  pcd.AddStep(func(pcd donothing.Process) {
    pcd.Name("retrieveBackupFile")
    pcd.Short("Retrieve the backup file")
  })
  pcd.AddStep(func(pcd donothing.Process) {
    pcd.Name("loadBackupData")
    pcd.Short("Load the backup data into the database")
  })
  pcd.AddStep(func(pcd donothing.Process) {
    pcd.Name("testRestoredData")
    pcd.Short("Check the restored data for consistency")
  })
}
```

When this code is run, the user will be prompted to follow the instructions step by step:

```markdown
# Restore a database backup

## Retrieve the backup file

Press `Enter` when done:

## Load the backup data into the database

Press `Enter` when done:

## Check the restored data for consistency

Press `Enter` when done:

Done!
```

Details for each step may be added in template files. For example, if we put the following markdown
into the file `templates/retrieveBackupFile.md`,

```markdown
- Log in to the AWS console
- Navigate to the `db_backups` S3 bucket
- Download the `.sql` file you wish to restore
```

then, when we run our code, the output will look like this:

```markdown
# Restore a database backup

## Retrieve the backup file

- Log in to the AWS console
- Navigate to the `db_backups` S3 bucket
- Download the `.sql` file you wish to restore

Press `Enter` when done:
```

The main idea behind `donothing` is that, when you're ready to automate a step instead of performing
it manually, you can just add a `Run()` function to the step. Continuing with the example above,
suppose we write a `retrieveBackupFile` function that downloads the latest database backup to our
working directory. We can then automate the first step of our procedure:

```go
// retrieveBackupFile downloads the latest backup file from the S3 bucket.
//
// It returns the path to the local file containing the backup data.
func downloadBackupFile() (string, error) {
  // ... use the AWS API to download the latest backup file ...
  return "backup.sql", nil
}

func main() {
  ...
  pcd.AddStep(func(pcd donothing.Process) {
    pcd.Name("retrieveBackupFile")
    pcd.Short("Retrieve the backup file")
    pcd.Run(func(pcd donothing.Process)) error {
      filename, err := downloadBackupFile()
      if err != nil {
        return err
      }
      pcd.OutputString("backupFilePath", filename)
    }
  })
  ...
}
```

Now when we run our script, it will automatically download the backup file from S3, allowing us to
move on to the second step immediately:

```markdown
# Restore a database backup

## Retrieve the backup file

    backupFilePath: ./backup.sql

## Load the backup data into the database

Press `Enter` when done:
```

This paradigm makes it easy to automate the restore procedure piece by piece.
```

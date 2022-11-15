# Database
database as a package which stores data in file system

## Example code

```
import (
	db "assignment/database"
)
var sourceDir = "/Users/XX/YY/ZZ/db"
client := db.SetupdbConnection(sourceDir)
err := client.Insert("name1", "akshay")
if err != nil {
	fmt.Println(err)
}
err = client.Update("name1", "mohan")
if err != nil {
	fmt.Println(err)
}
val, err := t.Get("name1")
fmt.Println(val, err)
```

## Features
DB will create an index file which will store all the indexes and also created a folder called "dbFiles" which will host all the files in which actual data will exist. Each file in that folder should not exceed size limit of 1mb(configurable). Once this size limit exceeds our code will create new file to store new data.
Index file will store data like this : {key : location of file}
- Get :
    - First checks key in the index file.
    - If key is not found return "Key not found"
    - Else 
        - Get the location for the db file for that key from index file
        - Then do a lookup in that db file and returns the value
- Insert :
    - Get the latest modified file from db folder
    - If size of that file is more than 1 mb:
        - Creates a new file and make the entry
    - Else make entry to the latest modified file
    - Make entry in index file
- Update :
    - Get location from index file
    - Update entry

## Optimizations

- For read queries we can further optimize the code by reading the db file chunks concurrently by using multiple go routines

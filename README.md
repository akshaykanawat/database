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
val, err := client.Get("name1")
fmt.Println(val, err)
```

## Features
The database will produce an index file to keep all the indexes and a folder called "dbFiles" to house all the files containing the actual data.  
There is a 1 MB size restriction on each file in that folder (configurable). Our code will generate a new file to hold fresh data once this size restriction is exceeded. Index file will store data like this : {key : location of file}
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

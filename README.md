# LotSize Exporter
## This Script was designed to generate lot size lists from Data available on NSE India's website
#### There is a file structure to follow: keep all the files in one dir, you can use sub dirs

The data comes as a zip download, you can extract using the following bash command

```bash
find . -name "*.zip" | xargs -P 5 -I fileName sh -c 'unzip -o -d "$(dirname "fileName")/$(basename -s .zip "fileName")" "fileName"'
```
```
Note: You may replace -P 5 with how many ever processes you are comfortable running on your own system.
```

The above command will unzip the archives into folders matching the archive name. Just remember to put all the downloaded zips in one location. 

you may also change the output dir to your liking. By default it will write to the same dir the code is in.

#### To Run:
```bash
# a simple
go run main.go
```

# asnlookup
CLI and Go package for fast, offline ASN lookups.

A level compressed trie in array representation is used for achieving very fast
lookups with a small memory footprint. The level compression is user-tunable 
between space-efficiency and time-efficiency. In LC-trie terms, the tuning
adjusts the fill factor of the redundancy-enabled level compression.

Due to the array-represented trie and binary marshaling, the inflation of a
pre-converted database can be measured in tens of milliseconds. In other words
the CLI tool can even be used for one-off lookups without any perceivable 
startup slowness.

```
time asnlookup --db ~/.asnlookup.db 8.8.8.8
15169

real    0m0,027s
user    0m0,025s
sys     0m0,018s
```

## Installation

Using prebuilt binaries:
```shell
curl -fsSL https://github.com/banviktor/asnlookup/releases/download/v0.1.1/asnlookup-linux-amd64-v0.1.1.tar.gz | sudo tar -zx -C /usr/local/bin
```

From source:
```shell
make
sudo make install
```

## Usage

### CLI

1. Download a fresh RIB dump, e.g. from http://archive.routeviews.org/:
    ```shell
    ./hack/pull_rib.sh
    ```
2. Convert it to `asnlookup`'s own format:
    ```shell
    bzcat rib.*.bz2 | asnlookup-utils convert --input - --output /path/to/my.db
    ```
3. Use it with `asnlookup`:
    ```shell
    asnlookup --db /path/to/my.db 8.8.8.8
    ```
   or using the `ASNLOOKUP_DB` environment variable:
    ```shell
    export ASNLOOKUP_DB=/path/to/my.db
    asnlookup 8.8.8.8
    ```

#### Batch lookups

You may also do batch lookups for IPs provided to standard input using the 
`--batch` flag:
```shell
echo -ne '1.1.1.1\n8.8.8.8\n' | asnlookup --db ~/.asnlookup.db --batch
13335
15169
```
If you have tons of IPs to check, this will be a lot faster than inflating the
multi-megabyte database each time `asnlookup` is invoked.

### Go package

1. Build a database

   * Manually:
     ```go
     builder := database.NewBuilder()

     _, prefix, _ := net.ParseCIDR("8.8.0.0/16")
     err := builder.InsertMapping(prefix, 420)
     if err != nil {
         panic(err)
     }
        
     db, err := builder.Build()
     if err != nil {
         panic(err)
     }
     ```
   
   * Using an MRT file:
     ```go
     mrtFile, err := os.OpenFile("/path/to/file.mrt", os.O_RDONLY, 0)
     if err != nil {
         panic(err)
     }
     defer mrtFile.Close() 

     builder := database.NewBuilder()
     if err = builder.ImportMRT(mrtFile); err != nil {
         panic(err)
     }
    
     db, err := builder.Build()
     if err != nil {
         panic(err)
     }
     ```

   * Using a marshaled database (see `asnlookup-utils convert`):
     ```go
     dbFile, err := os.OpenFile("/path/to/file.db", os.O_RDONLY, 0)
     if err != nil {
         panic(err)
     }
     defer dbFile.Close()
    
     db, err := database.NewFromDump(dbFile)
     if err != nil {
         panic(err)
     }
     ```
   
2. Look things up!
   ```go
   as, err := db.Lookup(net.ParseIP("8.8.8.8"))
   if err != nil {
       panic(err)
   }
   fmt.Println(as.Number)
   ```

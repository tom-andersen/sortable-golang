# Sortable Test written in Go

To compile, run:

```
./run.sh
```

An executable will be create. On Windows, this is *main.exe*:

```
$ ./main.exe

Please provide path to listings file.
  -listings string
        path to listings file
  -output string
        path to output file
  -products string
        path to products file
```

If you place the *products.txt* and *listings.txt* in local folder, then you can run:

```
$ ./main.exe -listings listings.txt -products products.txt -output results.txt

20196 listings read.
3445 listing with unknown manufucturer.
7290 matches written.
```

This will create a *results.txt* file.
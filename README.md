# gosort – Concurrent Chunk Sorting Assignment

**Student:** Ahmet Can Karayoluk**  
**Student ID:** 231ADB260**  
**Course:** DE0917(English)(1), 25/26-R**

---

## Overview

This program implements a concurrent integer sorting tool named **gosort**.  
It follows the required rules:

- Minimum 4 chunks  
- Otherwise: number of chunks = `ceil(sqrt(n))`  
- Chunks are nearly equal in size  
- Each chunk is sorted in a separate goroutine  
- Manual merge procedure (no flatten + global sort)  
- Modes: `-r`, `-i`, `-d`

---

## Usage

Run the program with:

go run .

shell
Kodu kopyala

Then choose one mode.

### Random Mode (`-r`)
go run . -r 20

mathematica
Kodu kopyala

### Input File Mode (`-i`)
go run . -i input.txt

shell
Kodu kopyala

### Directory Mode (`-d`)
Sorted files are written to:

incoming_sorted_ahmet_can_karayoluk_231ADB260

Kodu kopyala
go run . -d incoming

pgsql
Kodu kopyala

---

## Notes

- Sorting is performed concurrently.  
- Merging is implemented manually.  
- All modes match the assignment requirements.
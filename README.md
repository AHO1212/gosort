# gosort – Concurrent Chunk Sorting Assignment

**Student:** Ahmet Can Karayoluk  
**Student ID:** 231ADB260  
**Course:** DE0917(English)(1), 25/26-R  
**Assignment:** Concurrent Chunk Sorting (gosort)

---

## Overview

This program implements a concurrent integer sorting tool named **gosort**, following all requirements defined in the assignment:

- Minimum 4 chunks  
- Otherwise: number of chunks = `ceil(sqrt(n))`  
- Chunks must have almost equal size  
- Each chunk is sorted in its own goroutine  
- Merging is implemented manually (no flatten + global sort)  
- Three modes are supported: `-r`, `-i`, `-d`  

---

## Usage

### Random Mode (`-r`)
Generates N random integers (N ≥ 10) and sorts them concurrently.

gosort -r 20

### Input File Mode (`-i`)
Reads integers from a text file (one per line) and sorts them concurrently.

gosort -i input.txt

### Directory Mode (`-d`)
Processes all `.txt` files inside a directory.  
Each file is sorted independently and written to:

incoming_sorted_ahmet_can_karayoluk_231ADB260

Command:

gosort -d incoming

---

## Notes

- All chunk sorting is performed using goroutines.  
- Merging is done manually using multi-step k-way merging.  
- Input validation and error handling follow assignment requirements.  
- This implementation fully satisfies the requirements for modes `-r`, `-i`, and `-d`.  

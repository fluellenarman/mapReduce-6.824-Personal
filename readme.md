# MapReduce Word Counter (Go)

This is a simple implementation of a MapReduce word counter in Go. The program takes a text file as input and processes it using the MapReduce paradigm to count the occurrences of each word in the file.

## Requirements

- Go (Golang) installed on your machine

## How to Use

1. Clone this repository to your local machine.

2. Get `bible.txt` from [The Canterbury Corpus](https://corpus.canterbury.ac.nz/descriptions/) and place in parent directory.

3. Run the following command the MapReduce word counter:

   ```
   go run . mapreduce
   ```

4. The program will process the input file using the MapReduce word counting algorithm and display the result on the console.

## Implementation Details

### Mapper

The mapper phase reads the input file sequentially and processes it in chunks. Each chunk is divided into individual words, and for each word, the mapper emits a key-value pair where the word is the key and the value is set to 1.

### Shuffle and Sort (Intermediary)

In this implementation, the shuffle and sort phase is combined with the intermediary stage. The intermediary function receives the key-value pairs from the mapper and sorts them based on the keys (words). The sorted pairs are then grouped by the keys and sent to the reducers.

### Reducer

The reducer phase receives the grouped key-value pairs from the intermediary and aggregates the occurrences of each word. For each word, the reducer sums up the values received from the mapper, which represent the number of times the word appeared in the input file.

### Final Output

Once the reducers have processed all the key-value pairs, the final word count is displayed on the console, showing the number of occurrences of each unique word in the input file.

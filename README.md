# Min-𝜖 Cosine 𝑘-Nearest Neighbor Graph Construction

Given a set of objects 𝐷, for each object 𝑑𝑖 in 𝐷, find the 𝑘 most similar other objects 𝑑𝑗 with cosine similarity cos(𝑑𝑖, 𝑑𝑗) of at least 𝜖.

###Method Used:
IdxJoin (with minor improvements)

###Datasets: 
wiki1k.csr, wiki10.csr, wiki100.csr, wiki250.csr and sample.csr are sets of 1000, 10, 100, 250 and 3 documents,
respectively, randomly chosen from the English Wikipedia, downloaded in Oct,2014. Documents have been processed by stemming and stop-word removal, and are represented as term-frequency vectors. Documents have a minimum of 200 individual terms. The datasets are in CSR format, containing one document in each line of the text file.

###Instructions:

cd min-cosine-k-nearest/build

make

cd datasets

../findsim -eps 0.9 -k 10 wiki1.csr wiki1.nbrs.9.10.csr

The last command will invoke the ij mode (IdxJoin algorithm) in the findsim program with inputs 𝜖 = 0.9, 𝑘 = 10, and dataset wiki1.csr, creating the output wiki1.nbrs.9.10.csr.

Note that the output will be different for different values of 𝜖 and 𝑘. The output of IdxJoin are sparse vectors, sorted in decreasing similarity order, where the feature ID is the ID of
the neighboring document. IDs in input and output files are 1-indexed, i.e., they start at 1.

###Usage: 

findsim [options] input-file [output-file]

<input/output-file> should be in CSR.

Input is assumed to be a document term-frequency matrix.

If no <output-file> is specified, the output will not be saved. K-NNG output will be sparse vectors, sorted in decreasing similarity order.

####Options:

-k=int

Number of neighbors to return for each row in the Min-eps K-Nearest Neighbor Graph.

Default value is 10.

-eps=float

Minimum similarity for neighbors.

Default value is 0.5. Must be non-negative.

####Credits:
Inspired heavily by Prof. David C. Anastasiu

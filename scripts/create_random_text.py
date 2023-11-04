"""
Script to generate 1000 random documents with random words and add them to
the DocuStore database.

Used for testing purposes.
"""

import os
import random
import subprocess
from time import perf_counter

DOC_LENS = [100, 1000, 5000]


def main():
    vocab = open("./assets/words.txt").read().splitlines()
    for i in range(1000):
        doc_len = random.choice(DOC_LENS)
        doc = " ".join(random.choices(vocab, k=doc_len))
        with open(f"/tmp/text_{i}.txt", "w") as f:
            f.write(doc)
        init = perf_counter()
        subprocess.call(["./build/bin/DocuStore", "add", f"/tmp/text_{i}.txt"])
        print(f"{i+1} time: {1000 * (perf_counter() - init):.2f} ms")
        os.remove(f"/tmp/text_{i}.txt")


if __name__ == "__main__":
    main()

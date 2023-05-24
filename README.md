# DocuStore

A program to save any URL or raw text (with support for Markdown) and search for them later.

## Motivation

With browser favorites, you can only search webpage titles. The DocuStore search does a full text search on all saved webpages and returns the results ranked by similarity. Therefore, it is much easier to find that one webpage you are looking for.

Also, the support for Markdown allows to conveniently store your own documents to be searched in the future. Markdown content is rendered directly inside the app.

## Example screenshot

![Example screenshot](https://github.com/mathpn/DocuStore/raw/main/screenshots/sample_screenshot.png?raw=true)

## How does it work

DocuStore is written in Golang with a simple Javascript (Vue.js) frontend. Wails is used to package both the golang code and the frontend into a single nice binary file. Also thanks to Wails, the program is cross-platform and can be built for Windows, Linux and Mac.

### A summary of its inner parts

1. When a URL is provided, the raw HTML source code is parsed to extract relevant text information. When Markdown is provided, this step is skipped.
2. The raw text is tokenized and a data structure with token counts is created and persisted to disk using SQLite.
3. The new document is added to an inverted-index (stored in a self-balancing binary search tree).
4. When a search query is provided, all relevant documents are retrieved using the inverted-index.
5. A TF-IDF similarity is calculated and used to rank all documents, which are then returned in ranked order.

## CLI

You're not a GUI fan? That's OK. You can add URLs or Markdown files with the following syntax:

```bash
./DocuStore add <URL_OR_FILEPATH>
```

Later, you can query your stored documents with:

```bash
./DocuStore query <QUERY_STRING>
```

Where ./DocuStore is the path to the DocuStore binary.

## Requirements

If you want to build or develop locally, you need to have [Golang](https://go.dev/) and [Wails](https://wails.io/) installed.

## Building

To build a redistributable, production mode package, use `wails build`.

## License

BSD-3


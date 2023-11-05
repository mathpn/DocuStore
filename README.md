# DocuStore

A solution to store and retrieve URLs and Markdown documents.

## Why DocuStore?

I was tired of sifting through browser favorites, searching for that one elusive webpage. Traditional favorites only allow you to search by webpage titles, leaving you in the dark when you need to find specific information buried within a saved page. DocuStore performs a full-text search across all your saved web content and documents, presenting results ranked by relevance.

### Markdown Support

DocuStore also provides support for Markdown documents. You can conveniently store your own content for future reference. DocuStore also renders Markdown directly within the app.

## Search Demo

![GIF showing a search demo](https://github.com/mathpn/DocuStore/raw/main/assets/search_demo.gif?raw=true)

## Getting Started

- Linux

  Ensure you have `libgtk3` and `libwebkit` installed. Exact names depend on your Linux distro. Then, download the pre-compiled binary from the [releases](https://github.com/mathpn/DocuStore/releases/) section. For more information, check the [platform-specific dependencies](https://wails.io/docs/gettingstarted/installation#platform-specific-dependencies) section of the Wails documentation.

- Any platform

  If you want to build or develop locally, you need to have both [Go](https://go.dev/) and [Wails](https://wails.io/) installed.

## How does it work

DocuStore is written in Go with a simple Javascript (Vue.js) frontend, neatly packaged into a single binary file thanks to Wails. This means it's a cross-platform application that can be built for Linux, Mac and Windows.

### A summary of its inner parts

1. When you provide a URL, DocuStore parses the raw HTML source code to extract the most relevant text information. If you're adding Markdown, this step is skipped.
2. The raw text is tokenized and converted into a data structure with token counts, which is persisted to disk using SQLite.
3. Your new document is integrated into an inverted index.
4. When you run a search query, the inverted index is used to retrieve relevant documents.
5. Documents become [TF-IDF](https://en.wikipedia.org/wiki/Tf%E2%80%93idf) vectors, and they're ranked according to the cosine similarity to your query.

## Command Line Interface (CLI)

Not a fan of graphical interfaces? No problem. You can interact with DocuStore via the command line:

- Add URLs or Markdown files using the following syntax:

```bash
./DocuStore add <URL_OR_FILEPATH>
```

Later, you can query your stored documents using:

```bash
./DocuStore query <QUERY_STRING>
```

Where ./DocuStore is the path to the DocuStore binary.

## License

BSD-3

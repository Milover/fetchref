# fetchref

A simple, ~~hacky~~ command line utility for fetching article PDFs from
[Sci-Hub][Sci-Hub], book PDFs from [Libgen][Libgen] and formatted citations
from [CrossRef][CrossRef] from supplied DOIs and/or ISBNs.

Article download requires DOIs, while book downloads require ISBNs.
Citations, on the other hand, can be fetched both from DOIs and ISBNs,
although it should be noted that ISBNs are converted to DOIs
by querying [CrossRef][CrossRef] which sometimes fails, especially for
older books.

## TODO

- [ ] release stuff
    - [ ] setup GitHub actions/releases
    - [ ] publish on [pkg.go.dev](https://pkg.go.dev/)
- [ ] add more functionality for managing citations
- [ ] add better logging/instrumentation
- [ ] add unit/integration tests

[Sci-Hub]: https://sci-hub.se
[Libgen]: https://libgen.is
[CrossRef]: https://www.crossref.org

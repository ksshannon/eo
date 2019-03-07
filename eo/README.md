# eo

This repository houses executive order data from archive.gov, as well as the
federal register.  It provides a simple parser to extract the archive.gov data.
There are two sources utilized here, the second is federalregister.gov, but
that data source doesn't have various metadata for all orders (revokes, amends,
notes, etc.).  The federal register data can be accessed easily via a REST API,
the archive.gov data must be scraped or manually downloaded.  The
federalregister.gov data is updated locally as data/fr.json and imported after
the text EO data.  The internal representation of the data is the shorter of
the two (early EO data).  The data is ingested and can be exported to a csv,
yaml, or json file using cmd/export.go

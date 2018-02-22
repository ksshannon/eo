# eo

This repository houses executive order data from archive.gov, and provides a
parser to extract the data.  There are two sources utilized here, the second is
federalregister.gov, but that data source doesn't have various metadata for all
orders (revokes, amends, notes, etc.).  The federal register data can be
accessed easily via a REST API, the archive.gov data must be scraped or
manually downloaded.  Currently, updated data is copy/pasted into the
appropriate data/year.txt, and the parser chugs through those files.  Currently
the whole process goes:

0. Note the last recorded order in the data folder
1. Navigate to the archive.gov [disposition tables](https://www.archives.gov/federal-register/executive-orders/disposition)
2. Locate the last recorded order, and copy paste all following orders into the
   appropriate file
3. Build the export command `go build cmd/export.go` from the eo folder
4. Run the command `./export -f csv eo.csv` to generate a csv file of executive
   order information

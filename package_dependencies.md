# Dependency Table For: github.com/flowdev/spaghetti-analyzer

| | a n a l d a t a - T | d e p s - S | d o c - S | p a r s e - S | s t a t - S | t r e e - S | x / d i r s - T | x / p k g s - T |
| :- | :- | :- | :- | :- | :- | :- | :- | :- |
| **/** | **T** | **S** | **S** | **S** | **S** | **S** | **T** | **T** |
| deps | T | | | | | | | T |
| doc | T | | | | | | | |
| parse | | | | | | | | T |
| stat | T | | | | | | | |
| tree | | | | | | | T | T |

### Legend

* Rows - Importing packages
* Columns - Imported packages


#### Meaning Of Row And Row Header Formatting

* **Bold** - God package (can use all packages)
* `Code` - Database package (can only use tool and other database packages)
* _Italic_ - Tool package (foundational, no dependencies)
* No formatting - Standard package (can only use tool and database packages)


#### Meaning Of Letters In Table Columns

* G - God package (can use all packages)
* D - Database package (can only use tool and other database packages)
* T - Tool package (foundational, no dependencies)
* S - Standard package (can only use tool and database packages)

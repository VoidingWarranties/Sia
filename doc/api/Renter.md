Renter API
==========

This document contains detailed descriptions of the renter's API routes. For an
overview of the renter's API routes, see [API.md#renter](/doc/API.md#renter).  For
an overview of all API routes, see [API.md](/doc/API.md)

There may be functional API calls which are not documented. These are not
guaranteed to be supported beyond the current release, and should not be used
in production.

Overview
--------

The renter manages the user's files on the network. The renter's API endpoints
expose methods for managing files on the network and managing the renter's
allocated funds.

Index
-----

| Route                                                      | HTTP verb |
| ---------------------------------------------------------- | --------- |
| [/renter/allowance](#renterallowance-get)                  | GET       |
| [/renter/allowance](#renterallowance-post)                 | POST      |
| [/renter/downloads](#renterdownloads-get)                  | GET       |
| [/renter/files](#renterfiles-get)                          | GET       |
| [/renter/delete/___:siapath___](#renterdeletesiapath-post) | POST      |
| [/renter/download/:siapath___](#renterdownloadsiapath-get) | GET       |
| [/renter/rename/___:siapath___](#renterrenamesiapath-post) | POST      |
| [/renter/upload/___:siapath___](#renteruploadsiapath-post) | POST      |

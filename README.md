# things cloud sdk

[Things](https://culturedcode.com/things/) comes with a cloud based API, which can
be used to synchronize data between devices.
This is a golang SDK to interact with that API, opening the API so that you
can enhance your Things experience on iOS and Mac.

[![wercker status](https://app.wercker.com/status/ddec74f2f7406079026aa44e8a004a86/s/master "wercker status")](https://app.wercker.com/project/byKey/ddec74f2f7406079026aa44e8a004a86)

## TODO

- [x] Verify Credentials
- [x] Account Management
  - [x] Signup/ Confirmation
  - [x] Change Password
  - [x] Account Deletion
- [x] History management
  - [x] List Histories
  - [x] Create History 
  - [x] Delete History
  - [x] Sync History
  - [ ] Item Management
    - [x] read items 
    - [x] write items
    - recurring tasks
      - [x] neverending
      - [x] end on date
      - [x] end after n times
      - [ ] repeat after completion
      - [ ] reminders
      - [ ] deadlines
  - [x] State aggregation
    - [x] InMemory
    - [ ] Persistent

## Note

As there is no official API documentation available all requests need to be reverse engineered,
which takes some time. Feel free to contribute and improve & extend this implementation.

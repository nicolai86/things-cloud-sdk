# things cloud sdk

[Things](https://culturedcode.com/things/) comes with a cloud based API, which can
be used to synchronize data between devices.
This is a golang SDK to interact with that API, opening the API so that you
can enhance your Things experience on iOS and Mac.

[![wercker status](https://app.wercker.com/status/ddec74f2f7406079026aa44e8a004a86/s/master "wercker status")](https://app.wercker.com/project/byKey/ddec74f2f7406079026aa44e8a004a86)

## TODO

- [x] Verify Credentials
- [ ] History management
  - [ ] List Histories
  - [ ] Create History 
  - [ ] Delete History
- [ ] Item Management
  - [ ] Tasks
    - [ ] Create
    - [ ] Read
    - [ ] Update
    - [ ] Delete
    - [ ] Status (Created, Completed, Cancelled)
  - [ ] Checklists
    - [ ] Create
    - [ ] Read
    - [ ] Update
    - [ ] Delete
    - [ ] Status (Todo, Done)
  - [ ] Projects
    - [ ] Create
    - [ ] Read
    - [ ] Update
    - [ ] Delete
  - [ ] Relationships
    - [ ] Tasks -> Checklists
    - [ ] Projects -> Tasks

## Note

As there is no official API documentation available all requests need to be reverse engineered,
which takes some time. Feel free to contribute and improve & extend this implementation.

## ADDED Requirements

### Requirement: README SHALL contain a Version History section
README.md SHALL contain a second-level heading `## Version History` as a dedicated section for recording project version changes.

#### Scenario: Version History section exists
- **WHEN** reading README.md
- **THEN** the file SHALL contain a line matching `## Version History`

#### Scenario: Section position
- **WHEN** reading README.md
- **THEN** `## Version History` SHALL appear after all other existing sections

### Requirement: Version entries SHALL use table format
The Version History section SHALL contain a Markdown table with columns: Version, Date, Changes.

#### Scenario: Table header present
- **WHEN** reading the Version History section
- **THEN** the table SHALL have headers `| Version | Date | Changes |`

#### Scenario: At least one version entry
- **WHEN** reading the Version History section
- **THEN** the table SHALL contain at least one data row

### Requirement: Version entry format
Each version entry SHALL contain a semantic version number, an ISO 8601 date (YYYY-MM-DD), and a brief change description.

#### Scenario: Valid version number format
- **WHEN** reading a version entry
- **THEN** the Version column SHALL match the pattern `vX.Y.Z` where X, Y, Z are non-negative integers

#### Scenario: Valid date format
- **WHEN** reading a version entry
- **THEN** the Date column SHALL match the pattern `YYYY-MM-DD`

#### Scenario: Non-empty changes description
- **WHEN** reading a version entry
- **THEN** the Changes column SHALL be non-empty

### Requirement: Version entries SHALL be in reverse chronological order
Entries in the version history table SHALL be ordered with the most recent version first.

#### Scenario: Descending order by version
- **WHEN** reading the version history table with multiple entries
- **THEN** version numbers SHALL decrease from top to bottom

### Requirement: Edge case - empty changes column
The system SHALL NOT allow empty Changes descriptions in version entries.

#### Scenario: Reject empty changes
- **WHEN** a version entry has an empty Changes column
- **THEN** the acceptance test SHALL fail

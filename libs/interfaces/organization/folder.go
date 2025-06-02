package organization

type Folder struct {
    Id          string
    Name        string
    Description *string
    Icon        *string
    CreatedAt   string
    UpdatedAt   string
    ParentId    *string
    Members     []string
    CreatedBy   string
}

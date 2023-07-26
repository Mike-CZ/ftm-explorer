package resolvers

// NumberOfAccounts returns the number of accounts.
func (rs *RootResolver) NumberOfAccounts() int32 {
	return int32(rs.repository.GetNumberOfAccounts())
}

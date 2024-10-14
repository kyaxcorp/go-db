package dbresolver

func (r *resolver) CanProcessWriteOp() bool {
	if r.resolverStatus.Get() == ResolverReadyToProcess {
		return true
	}
	return false
}

func (r *resolver) CanProcessReadOp() bool {
	switch r.resolverStatus.Get() {
	case ResolverReadOnly:
		return true
	case ResolverReadyToProcess:
		return true
	}
	return false
}

func (r *resolver) CanProcessTransactionOp() bool {
	return r.CanProcessWriteOp()
}

func (r *resolver) CanProcessTheseTables(tables []string) bool {
	// TODO: we should know if tables are accessible
	return true
}

package trust

import (
	"fmt"

	"github.com/theupdateframework/notary/client"
	"github.com/theupdateframework/notary/trustpinning"
	"github.com/theupdateframework/notary/tuf/data"
)

// SignAndPublish signs an artifact, then publishes the metadata to a trust server
func SignAndPublish(trustDir, trustServer, gun, tlscacert, rootKey string) error {
	if err := ensureTrustDir(trustDir); err != nil {
		return fmt.Errorf("cannot ensure trust directory: %v", err)
	}

	transport, err := makeTransport(trustServer, gun, tlscacert)
	if err != nil {
		return fmt.Errorf("cannot make transport: %v", err)
	}

	repo, err := client.NewFileCachedRepository(
		trustDir,
		data.GUN(gun),
		trustServer,
		transport,
		getPassphraseRetriever(),
		trustpinning.TrustPinConfig{},
	)

	if err != nil {
		return fmt.Errorf("cannot create new file cached repository: %v", err)
	}

	err = clearChangeList(repo)
	if err != nil {
		return fmt.Errorf("cannot clear change list: %v", err)
	}

	rootKeyIDs, err := importRootKey(rootKey, repo, getPassphraseRetriever())
	if err != nil {
		return err
	}

	if err = repo.Initialize(rootKeyIDs); err != nil {
		return fmt.Errorf("ERROR: %v", err)
	}

	return nil
}

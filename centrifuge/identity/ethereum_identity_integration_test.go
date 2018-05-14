// +build ethereum

package identity_test

import (
	"testing"
	"github.com/spf13/viper"
	"os"
	"github.com/stretchr/testify/assert"
	"github.com/CentrifugeInc/go-centrifuge/centrifuge/tools"
	"github.com/CentrifugeInc/go-centrifuge/centrifuge/identity"
	"fmt"
	"github.com/CentrifugeInc/go-centrifuge/centrifuge/keytools"
)

func TestMain(m *testing.M) {
	//for now set up the env vars manually in integration test
	//TODO move to generalized config once it is available
	viper.BindEnv("ethereum.gethSocket", "CENT_ETHEREUM_GETH_SOCKET")
	viper.BindEnv("ethereum.gasLimit", "CENT_ETHEREUM_GASLIMIT")
	viper.BindEnv("ethereum.gasPrice", "CENT_ETHEREUM_GASPRICE")
	viper.BindEnv("ethereum.contextWaitTimeout", "CENT_ETHEREUM_CONTEXTWAITTIMEOUT")
	viper.BindEnv("ethereum.accounts.main.password", "CENT_ETHEREUM_ACCOUNTS_MAIN_PASSWORD")
	viper.BindEnv("ethereum.accounts.main.key", "CENT_ETHEREUM_ACCOUNTS_MAIN_KEY")
	viper.BindEnv("identity.ethereum.identityFactoryAddress", "CENT_IDENTITY_ETHEREUM_IDENTITYFACTORYADDRESS")
	viper.BindEnv("identity.ethereum.identityRegistryAddress", "CENT_IDENTITY_ETHEREUM_IDENTITYREGISTRYADDRESS")
	viper.Set("identity.ethereum.enabled", "true")
	viper.Set("keys.signing.publicKey", "../../resources/signingKey.pub")
	viper.Set("keys.signing.privateKey", "../../resources/signingKey.key")

	result := m.Run()
	os.Exit(result)
}

func TestCreateAndResolveIdentity_Integration(t *testing.T) {
	centrifugeId := tools.RandomString32()
	nodePeerId := tools.RandomByte32()

	identityObj := identity.NewEthereumIdentity()
	identityObj.CentrifugeId = centrifugeId
	identityObj.Keys[1] = append(identityObj.Keys[1], identity.EthereumIdentityKey{nodePeerId})

	confirmations := make(chan *identity.EthereumIdentity, 1)

	err := identity.CreateEthereumIdentity(identityObj, confirmations)
	assert.Nil(t, err, "should not error out when creating identity")

	registeredIdentity := <-confirmations
	assert.Equal(t, centrifugeId, registeredIdentity.CentrifugeId, "Resulting Identity should have the same ID as the input")

	id, err := identity.ResolveP2PEthereumIdentityForId(centrifugeId)
	assert.Nil(t, err, "should not error out when resolving identity")
	assert.Equal(t, centrifugeId, id.CentrifugeId, "CentrifugeId Should match provided one")
	assert.Equal(t, 0, len(id.Keys), "Identity Should have empty map of keys")
}

func TestCheckIdentityExists_Integration(t *testing.T) {
	centrifugeId := tools.RandomString32()
	nodePeerId := tools.RandomByte32()

	identityObj := identity.NewEthereumIdentity()
	identityObj.CentrifugeId = centrifugeId
	identityObj.Keys[1] = append(identityObj.Keys[1], identity.EthereumIdentityKey{nodePeerId})

	confirmations := make(chan *identity.EthereumIdentity, 1)

	err := identity.CreateEthereumIdentity(identityObj, confirmations)
	assert.Nil(t, err, "should not error out when creating identity")

	registeredIdentity := <-confirmations
	assert.Equal(t, centrifugeId, registeredIdentity.CentrifugeId, "Resulting Identity should have the same ID as the input")

	exists, err := identityObj.CheckIdentityExists()
	assert.Nil(t, err, "should not error out when looking for correct identity")
	assert.Equal(t, true, exists, "Identity Should Exists")

	wrongCentrifugeId := tools.RandomString32()
	identityObj.CentrifugeId = wrongCentrifugeId
	exists, err = identityObj.CheckIdentityExists()
	assert.Nil(t, err, "should not error out when missing identity")
	assert.Equal(t, false, exists, "Identity Should Exists")
}

func TestCreateIdentityAndAddKey_Integration(t *testing.T) {
	centrifugeId := tools.RandomString32()
	nodePeerId := tools.RandomByte32()

	identityObj := identity.NewEthereumIdentity()
	identityObj.CentrifugeId = centrifugeId
	identityObj.Keys[1] = append(identityObj.Keys[1], identity.EthereumIdentityKey{nodePeerId})

	confirmations := make(chan *identity.EthereumIdentity, 1)

	err := identity.CreateEthereumIdentity(identityObj, confirmations)
	assert.Nil(t, err, "should not error out when creating identity")

	registeredIdentity := <-confirmations
	assert.Equal(t, centrifugeId, registeredIdentity.CentrifugeId, "Resulting Identity should have the same ID as the input")

	id, err := identity.ResolveP2PEthereumIdentityForId(centrifugeId)

	assert.Nil(t, err, "should not error out when resolving identity")
	assert.Equal(t, centrifugeId, id.CentrifugeId, "CentrifugeId Should match provided one")
	assert.Equal(t, 0, len(id.Keys), "Identity Should have empty map of keys")

	err = identityObj.AddKeyToIdentity(1, confirmations)
	assert.Nil(t, err, "should not error out when adding key to identity")

	receivedIdentity := <-confirmations
	assert.Equal(t, centrifugeId, receivedIdentity.CentrifugeId, "Resulting Identity should have the same ID as the input")
	assert.Equal(t, 1, len(receivedIdentity.Keys), "Resulting Identity Key Map should have expected length")
	assert.Equal(t,1, len(receivedIdentity.Keys[1]), "Resulting Identity Key Type list should have expected length")
	assert.Equal(t, identityObj.Keys[1][0].Key, receivedIdentity.Keys[1][0].Key, "Resulting Identity Key should match the one requested")

	// Double check that Key Exists in Identity
	id, err = identity.ResolveP2PEthereumIdentityForId(centrifugeId)

	assert.Nil(t, err, "should not error out when resolving identity")
	assert.Equal(t, centrifugeId, id.CentrifugeId, "CentrifugeId Should match provided one")
	assert.Equal(t, 1, len(id.Keys), "Identity Should have empty map of keys")
	assert.Equal(t, identityObj.Keys[1][0].Key, id.Keys[1][0].Key, "Resulting Identity Key should match the one requested")
}

func TestManageIdentity(t *testing.T) {
	// Creation Succeeds
	centrifugeId := tools.RandomString32()
	publicKey, _ := keytools.GetSigningKeyPairFromConfig()
	b32Key := tools.ConvertByteArrayToByte32(publicKey)

	err := identity.CreateEthereumIdentityFromApi(centrifugeId, b32Key)
	assert.Nil(t, err, "should not error out upon identity creation")

	// AddKey Again Fails as Key already exists
	err = identity.AddKeyToIdentityFromApi(centrifugeId, identity.KEY_TYPE_PEERID, b32Key)
	assert.EqualError(t, err, "Key trying to be added already exists as latest. Skipping Update.", "should error out upon double key addition")

	// Creation fails as CentId already exists
	err = identity.CreateEthereumIdentityFromApi(centrifugeId, b32Key)
	assert.EqualError(t, err, fmt.Sprintf("ACTION [%v] but identity exists [%v]", identity.ACTION_CREATE, true), "should error out if ID already exists")

	// Adding Key fails as CentId does not exist
	centrifugeId = tools.RandomString32()
	err = identity.AddKeyToIdentityFromApi(centrifugeId, identity.KEY_TYPE_PEERID, b32Key)
	assert.EqualError(t, err, fmt.Sprintf("ACTION [%v] but identity exists [%v]", identity.ACTION_ADDKEY, false), "should error if ID does not exist")
}

func TestCreateAndResolveIdentity_Integration_Concurrent(t *testing.T) {
	var submittedIds [5]string

	howMany := cap(submittedIds)
	confirmations := make(chan *identity.EthereumIdentity, howMany)

	for ix := 0; ix < howMany; ix++ {
		centId := tools.RandomString32()
		identityObj := identity.NewEthereumIdentity()
		identityObj.CentrifugeId = centId
		submittedIds[ix] = centId
		err := identity.CreateEthereumIdentity(identityObj, confirmations)
		assert.Nil(t, err, "should not error out upon identity creation")
	}

	for ix := 0; ix < howMany; ix++ {
		singleIdentity := <-confirmations
		id, err := identity.ResolveP2PEthereumIdentityForId(singleIdentity.CentrifugeId)
		assert.Nil(t, err, "should not error out upon identity resolution")
		assert.Contains(t, submittedIds, id.CentrifugeId , "Should have the ID that was passed into create function [%v]", id.CentrifugeId)
	}
}
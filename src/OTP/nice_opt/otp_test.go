package otp

import (
	"github.com/stretchr/testify/require"

	"testing"
)

func TestKeyAllThere(t *testing.T) {
	k, err := NewKeyFromURL(`otpauth://totp/Example:alice@google.com?secret=JBSWY3DPEHPK3PXP&issuer=Example`)
	require.NoError(t, err, "failed to parse url")
	require.Equal(t, "totp", k.Type(), "Extracting Type")
	require.Equal(t, "Example", k.Issuer(), "Extracting Issuer")
	require.Equal(t, "alice@google.com", k.AccountName(), "Extracting Account Name")
	require.Equal(t, "JBSWY3DPEHPK3PXP", k.Secret(), "Extracting Secret")
}

func TestKeyIssuerOnlyInPath(t *testing.T) {
	k, err := NewKeyFromURL(`otpauth://totp/Example:alice@google.com?secret=JBSWY3DPEHPK3PXP`)
	require.NoError(t, err, "failed to parse url")
	require.Equal(t, "Example", k.Issuer(), "Extracting Issuer")
	require.Equal(t, "alice@google.com", k.AccountName(), "Extracting Account Name")
}

func TestKeyNoIssuer(t *testing.T) {
	k, err := NewKeyFromURL(`otpauth://totp/alice@google.com?secret=JBSWY3DPEHPK3PXP`)
	require.NoError(t, err, "failed to parse url")
	require.Equal(t, "", k.Issuer(), "Extracting Issuer")
	require.Equal(t, "alice@google.com", k.AccountName(), "Extracting Account Name")
}

func TestKeyWithNewLine(t *testing.T) {
	w, err := NewKeyFromURL(`otpauth://totp/Example:alice@google.com?secret=JBSWY3DPEHPK3PXP
`)
	require.NoError(t, err)
	sec := w.Secret()
	require.Equal(t, "JBSWY3DPEHPK3PXP", sec)
}

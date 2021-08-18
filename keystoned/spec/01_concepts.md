<!--
order: 1
-->

# Introduction

Keystone is a set of components aimed at making it easier for
non-crypto-native users to participate in blockchain
transactions. Keystone builds on the concepts introduced in the groups
module, using Cosmos groups to allow signature of transactions by
multiple keys wielded on behalf of an individual user, both by the
user (via a "user-agent") and by other agents acting on behalf of the
user.

# Concepts

## Key group

A key group is a set of cryptographic keys which are wielded on behalf
of an individual Cosmos "user" - a natural person. Keys may be stored
on user devices and wielded directly by a natural person, or they may
be stored, and wielded on behalf of the user by one or more
servers. These servers must be _trusted_ by the user in some capacity,
in order to make this possible.

Multiple keys in the group may be involved in signing
transactions. For example, a user may elect to sign using a key
resident on their device, and also with a key maintained by their
Keystone server. The user may have multiple devices with keys, and say
that two of their device keys and a server key must be involved in
signing. 

## Keystone Server

The Keystone server maintains one or more keys on behalf of individual
natural people. Keys will not be stored on the machine/container
hosting the Keystone server system, but via mechanisms such as
cloud-based keystores (AWS CloudHSM) or hardware systems such as HSMs.

## Creating a new key group

When a register request is received, the following steps are performed:

  1. If the user has not passed a public key and/or public key in the
     request (either does not have, or wishes not to reveal, their
     existing Cosmos address) then a new keypair is created in an
     HSM. The public key is used to create a new Cosmos address for
     the user.
  2. An admin group is created, containing two members. One of those
     is the address representing the user. The other is a (already
     existing) group address for the Keystone servers. Both members
     are given the same weight. 
  3. A new key group is created, containing two members. As for the
     admin group, the same two members will be present in this group -
     one address representng the user agent (device) and the other
     representing Keystone server group. Its creator will be given as
     the admin group address created in step 2. above. Both members
     are given the same weight. 
	 
## Signing transactions

A Keystone user requests signing of transactions by sending a Sign
request to the Keystone server. A signing request must always be
initiated by the user, even if it is the Keystone server that performs
the actual cryptographic signing on behalf of the user.

## Adding or removing keys to the group

A set of keys is represented by a group with multiple members
(keys). A group administrator must make any changes to the membership
of the group. Since, initially, both at least the user, and the
Keystone server may add new members, either the user or the Keystone
administrators may make such changes.

## Fees

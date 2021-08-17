<!--
order: 1
-->

# Introduction

Keystone is a set of components aimed at making it easier for
non-crypto-native users to participate in blockchain
transactions. Keystone builds on the concepts introduced in the groups
module, using Cosmos groups to allow signature of transactions by
multiple keys wielded on behalf of an individual user. 

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

## Signing transactions

## Adding or removing keys to the group

## Fees

# slsa-l3-provenance-demo

A tiny Go CLI whose release binaries are built on GitHub Actions with
[SLSA v1.0 Build Level 3](https://slsa.dev/spec/v1.0/levels#build-l3) provenance, using
GitHub's recommended method: [artifact attestations](https://docs.github.com/en/actions/concepts/security/artifact-attestations)
(`actions/attest`) generated inside a [reusable workflow](https://docs.github.com/en/actions/how-tos/secure-your-work/use-artifact-attestations/increase-security-rating).

```console
$ slsa-l3-provenance-demo --name Ada
Hello, Ada!
$ slsa-l3-provenance-demo --version
v0.1.0
```

## How the provenance is produced

| Workflow | Trigger | Role |
| --- | --- | --- |
| [`ci.yml`](.github/workflows/ci.yml) | Push to `main`, pull requests | gofmt, vet, test, build. No attestations. |
| [`release.yml`](.github/workflows/release.yml) | Push of a `v*` tag | Caller: runs the trusted builder, checks its output, publishes the GitHub Release. |
| [`build-attest.yml`](.github/workflows/build-attest.yml) | `workflow_call` | Trusted builder: builds 6 binaries and `SHA256SUMS`, then signs SLSA provenance for them. |

Artifact attestations on their own give SLSA Build L2. Build L3 also needs provenance that the
build itself can't forge. GitHub's documented way to get there is to build and sign in a reusable
workflow, then verify which workflow signed. The Sigstore certificate in each attestation is bound to
the identity of the reusable workflow (`.../build-attest.yml@refs/tags/<tag>`), not the caller's. So
`gh attestation verify --signer-workflow` rejects anything signed by another workflow.

The trusted builder is hardened further:

- It takes no inputs, so the calling workflow can't change what is built or attested.
- The tag is validated as semver inside the builder. Untrusted values reach scripts only through
  environment variables, never `${{ }}` interpolation.
- Build and test code runs in a job without `id-token: write`, so it can't request a signing
  certificate. A separate job does the signing.
- The signing job attests digests passed through job outputs, not files from a shared artifact.
  Nothing else in the run can swap what gets signed.
- Releases build on GitHub-hosted runners with no caches restored. Every action is pinned to a
  full commit SHA.
- Before publishing, the release job checks that the files match the attested digests and runs
  `gh attestation verify` on each one.

The trusted builder lives in this repository, so anyone who can push a tag here can also change
it. For organisation-wide use, GitHub recommends a shared builder repository verified with
`--signer-repo`.

### Transparency log

This repository is public, so `actions/attest` signs with the
[Sigstore public-good instance](https://openssf.org/blog/2023/10/03/running-sigstore-as-a-managed-service-a-tour-of-sigstores-public-good-instance/):

- The short-lived signing certificate is logged in Fulcio's certificate transparency log.
- The signature is recorded in the public [Rekor](https://docs.sigstore.dev/logging/overview/)
  transparency log.

Both records are permanent and public. They include the repository name, workflow path, commit SHA
and tag. The attestation is also stored in this repository's attestation API.

## Cutting a release

```sh
git tag -a v0.1.0 -m "v0.1.0"
git push origin v0.1.0
```

If a release fails, fix it on `main` and push a new tag rather than moving an existing one.

## Verifying a release

You need a recent [GitHub CLI](https://cli.github.com/); this was tested with v2.105.0.

```sh
REPO=Siddharth-Srinivas/slsa-l3-provenance-demo
TAG=v0.1.0
gh release download "$TAG" --repo "$REPO" --dir "demo-$TAG" && cd "demo-$TAG"

gh attestation verify "slsa-l3-provenance-demo_${TAG}_darwin_arm64" \
  --repo "$REPO" \
  --signer-workflow "$REPO/.github/workflows/build-attest.yml" \
  --source-ref "refs/tags/$TAG" \
  --deny-self-hosted-runners
```

| Flag | Enforces |
| --- | --- |
| `--repo` | The artifact was built from this repository. |
| `--signer-workflow` | The provenance was signed by the trusted builder. This is the Build L3 check. |
| `--source-ref` | The build ran for this release tag. |
| `--deny-self-hosted-runners` | The build ran on a GitHub-hosted runner. |

`SHA256SUMS` is attested too. You can verify it the same way and then check every file against it:

```sh
gh attestation verify SHA256SUMS --repo "$REPO" \
  --signer-workflow "$REPO/.github/workflows/build-attest.yml" \
  --source-ref "refs/tags/$TAG" --deny-self-hosted-runners
shasum -a 256 -c SHA256SUMS   # or: sha256sum -c SHA256SUMS
```

To verify against the Sigstore bundle attached to the release instead of GitHub's attestation API,
add `--bundle "slsa-l3-provenance-demo_${TAG}.sigstore.json"`. For fully offline verification, first
save the trusted roots while online with `gh attestation trusted-root > trusted_root.jsonl`, then add
`--custom-trusted-root trusted_root.jsonl`.

With `--bundle`, `gh` checks `--repo` only against the name recorded in the certificate when the
artifact was signed. Drop `--bundle` to confirm the repository still has that name.

## Development

```sh
go test ./...
go run . --name you
```

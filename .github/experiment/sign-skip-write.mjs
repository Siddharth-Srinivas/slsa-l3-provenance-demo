// Signs SLSA provenance for one file with the GitHub Sigstore instance and
// writes the bundle to disk only: no Rekor upload (GitHub's instance uses an
// RFC 3161 timestamp instead) and no write to the GitHub attestation API.
import { createHash } from 'node:crypto'
import { readFileSync, writeFileSync } from 'node:fs'
import { basename } from 'node:path'
import { attest, buildSLSAProvenancePredicate } from '@actions/attest'

const [subjectPath, bundlePath] = process.argv.slice(2)
const digest = createHash('sha256').update(readFileSync(subjectPath)).digest('hex')
const predicate = await buildSLSAProvenancePredicate()

const attestation = await attest({
  subjects: [{ name: basename(subjectPath), digest: { sha256: digest } }],
  predicateType: predicate.type,
  predicate: predicate.params,
  sigstore: 'github',
  skipWrite: true,
  token: process.env.GITHUB_TOKEN
})

writeFileSync(bundlePath, JSON.stringify(attestation.bundle))
console.log(`sha256=${digest}`)
console.log(`tlogID=${attestation.tlogID ?? 'none'}`)
console.log(`attestationID=${attestation.attestationID ?? 'none (not written to GitHub)'}`)

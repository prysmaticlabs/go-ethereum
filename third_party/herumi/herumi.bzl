load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

"""
Herumi's BLS library for go depends on
- herumi/mcl
- herumi/bls
- herumi/bls-eth-go-binary
"""

def bls_dependencies():
    _maybe(
        http_archive,
        name = "herumi_bls_eth_go_binary",
        strip_prefix = "bls-eth-go-binary-da18d415993a059052dfed16711f2b3bd03c34b8",
        urls = [
            "https://github.com/herumi/bls-eth-go-binary/archive/da18d415993a059052dfed16711f2b3bd03c34b8.tar.gz",
        ],
        sha256 = "69080ca634f8aaeb0950e19db218811f4bb920a054232e147669ea574ba11ef0",
        build_file = "@prysm//third_party/herumi:bls_eth_go_binary.BUILD",
    )
    _maybe(
        http_archive,
        name = "herumi_mcl",
        strip_prefix = "mcl-1b043ade54bf7e30b8edc29eb01410746ba92d3d",
        urls = [
            "https://github.com/herumi/mcl/archive/1b043ade54bf7e30b8edc29eb01410746ba92d3d.tar.gz",
        ],
        sha256 = "306bf22b747db174390bbe43de503131b0b5b75bbe586d44f3465c16bda8d28a",
        build_file = "@prysm//third_party/herumi:mcl.BUILD",
    )
    _maybe(
        http_archive,
        name = "herumi_bls",
        strip_prefix = "bls-989e28ede489e5f0e50cfc87e3fd8a8767155b9f",
        urls = [
            "https://github.com/herumi/bls/archive/989e28ede489e5f0e50cfc87e3fd8a8767155b9f.tar.gz",
        ],
        sha256 = "14b441cc66ca7e6c4e0542dcfc6d9f83f4472f0e7a43efaa1d3ea93e2e2b7491",
        build_file = "@prysm//third_party/herumi:bls.BUILD",
    )

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

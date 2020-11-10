# Generated file, do not modify by hand
# Generated by 'rbe_ubuntu_clang_gen' rbe_autoconfig rule
"""Definitions to be used in rbe_repo attr of an rbe_autoconf rule  """
toolchain_config_spec0 = struct(config_repos = ["prysm_toolchains"], create_cc_configs = True, create_java_configs = True, env = {"BAZEL_COMPILER": "clang", "BAZEL_LINKLIBS": "-l%:libstdc++.a", "BAZEL_LINKOPTS": "-lm:-static-libgcc", "BAZEL_USE_LLVM_NATIVE_COVERAGE": "1", "GCOV": "llvm-profdata", "CC": "clang", "CXX": "clang++"}, java_home = "/usr/lib/jvm/java-8-openjdk-amd64", name = "clang")
_TOOLCHAIN_CONFIG_SPECS = [toolchain_config_spec0]
_BAZEL_TO_CONFIG_SPEC_NAMES = {"3.7.0": ["clang"]}
LATEST = "sha256:d5fa14154811dff0886e4c808dc15f18c4bb8545a1ef3c53805a0db13564bdad"
CONTAINER_TO_CONFIG_SPEC_NAMES = {"sha256:d5fa14154811dff0886e4c808dc15f18c4bb8545a1ef3c53805a0db13564bdad": ["clang"]}
_DEFAULT_TOOLCHAIN_CONFIG_SPEC = toolchain_config_spec0
TOOLCHAIN_CONFIG_AUTOGEN_SPEC = struct(
        bazel_to_config_spec_names_map = _BAZEL_TO_CONFIG_SPEC_NAMES,
        container_to_config_spec_names_map = CONTAINER_TO_CONFIG_SPEC_NAMES,
        default_toolchain_config_spec = _DEFAULT_TOOLCHAIN_CONFIG_SPEC,
        latest_container = LATEST,
        toolchain_config_specs = _TOOLCHAIN_CONFIG_SPECS,
    )
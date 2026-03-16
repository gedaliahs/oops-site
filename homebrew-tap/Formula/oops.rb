class Oops < Formula
  desc "Undo for your terminal — backs up files before destructive commands"
  homepage "https://oops-cli.com"
  version "0.1.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://oops-cli.com/releases/oops_darwin_arm64.tar.gz"
      sha256 "c42817f6f556fa26fa60f2f875eb2dc9c4f1dc5bc9efdf6b53c0c021aed75c6a"
    else
      url "https://oops-cli.com/releases/oops_darwin_amd64.tar.gz"
      sha256 "299a10ee1bcb4232a65eb7df9dff527401a5a7777872a2740b59e4592a4e0083"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://oops-cli.com/releases/oops_linux_arm64.tar.gz"
      sha256 "dc07c92bfdee7d8d00d707e43f509cdd10f46c9774ffd66364b603011ad3d13b"
    else
      url "https://oops-cli.com/releases/oops_linux_amd64.tar.gz"
      sha256 "d5a6f7b9f520af51ff6f70b71dc05234de986608037fcecaf14a85fa700f8939"
    end
  end

  def install
    bin.install "oops"
  end

  def caveats
    <<~EOS
      Add this to your shell config:

        # zsh (~/.zshrc)
        eval "$(oops init zsh)"

        # bash (~/.bashrc)
        eval "$(oops init bash)"

        # fish (~/.config/fish/config.fish)
        oops init fish | source

      Then restart your shell.
    EOS
  end

  test do
    assert_match "oops", shell_output("#{bin}/oops --help")
  end
end

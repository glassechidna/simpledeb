class Simpledeb < Formula
  desc "simpledeb aims to be the simplest way to create an apt repo from a collection of .deb files."
  homepage "https://github.com/aidansteele/simpledeb"
  url "https://github.com/aidansteele/simpledeb/releases/download/v0.1.0/simpledeb_0.1.0_Darwin_x86_64.tar.gz"
  version "0.1.0"
  sha256 "562ae05827df20d2a422ffb9297bcd42df44dbc18fc2bf00c3d8a7c2a6a7801c"

  def install
    bin.install "simpledeb"
  end
end

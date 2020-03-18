function quit() {
    echo $1
    exit 1
}

echo "Checking gomu installation..."
which gomu > /dev/null && version=$(gomu version) && echo "Found previous installation: $version" || echo "No previous versions found"

echo "Checking repo version..."
which git-tagger > /dev/null || (echo "Installing git-tagger..." && go install github.com/hatchify/git-tagger)

tag=$(git-tagger -action=get) && echo "Current repo version: $tag"
if [[ $tag == $version ]]; then
    echo "Already running the latest version of gomu!"
    exit 0
fi

echo "Installing gomu $tag..."
go install -i -v -ldflags="-X main.version=$tag" -trimpath || quit "Failed to build gomu :("

which gomu > /dev/null && version=$(gomu version)
if [[ $tag == $version ]]; then
    echo "Updated gomu to version $version!"
else 
    echo "Failed to update gomu :("
fi


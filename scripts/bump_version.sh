#!/bin/bash

set -e  # Exit on error

# Skip if this is a version bump commit or has [skip version] tag
if git log -1 --pretty=%B | grep -q "\[skip version\]\|Update README for version v"; then
    echo "Skipping version bump for version-related commit"
    exit 0
fi

# Get the latest tag from Git
current_version=$(git describe --tags --abbrev=0)
IFS='.' read -r major minor patch <<< "${current_version#v}"

# Increment the patch version
if [ "$patch" -lt 20 ]; then
    patch=$((patch + 1))
else
    patch=0
    minor=$((minor + 1))
fi

# If minor reaches 20, reset it and increment major
if [ "$minor" -ge 20 ]; then
    minor=0
    major=$((major + 1))
fi

# Create the new version string
new_version="$major.$minor.$patch"

# Check if tag already exists
if git rev-parse "v$new_version" >/dev/null 2>&1; then
    echo "Error: Tag v$new_version already exists. Please resolve the conflict manually."
    exit 1
fi

# Tag the new version
echo "Creating new tag v$new_version..."
if ! git tag "v$new_version"; then
    echo "Error: Failed to create tag v$new_version"
    exit 1
fi

if ! git push origin "v$new_version"; then
    echo "Error: Failed to push tag v$new_version"
    git tag -d "v$new_version"  # Clean up local tag if push fails
    exit 1
fi

# Check if README.md needs updating
if ! grep -q "go install github.com/tedfulk/goatmeal@v$new_version" README.md; then
    # Update the README.md file with the new version
    sed -i.bak "s|go install github.com/tedfulk/goatmeal@.*|go install github.com/tedfulk/goatmeal@v$new_version|" README.md
    rm README.md.bak
    echo "README.md updated with version v$new_version."
    # Stage the README change
    git add README.md
    # Commit the README change without triggering another version bump
    git commit -m "[skip version] Update README for version v$new_version"
    git push origin
else
    echo "README.md already up to date."
fi

echo "Successfully bumped version to v$new_version"
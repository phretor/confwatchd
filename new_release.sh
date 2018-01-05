#!/bin/bash
# nothing to see here, just a utility i use to create new releases ^_^

CURRENT_VERSION=$(cat config/version.go | grep APP_VERSION | cut -d '"' -f 2)
TO_UPDATE=(
    config/version.go
)

echo -n "Current version is $CURRENT_VERSION, select new version: "
read NEW_VERSION
echo "Creating version $NEW_VERSION ...\n"

for file in "${TO_UPDATE[@]}"
do
    echo "Patching $file ..."
    sed -i "s/$CURRENT_VERSION/$NEW_VERSION/g" $file
    git add $file
done

git commit -m "Releasing v$NEW_VERSION"
git push

git tag -a v$NEW_VERSION -m "Release v$NEW_VERSION"
git push origin v$NEW_VERSION

echo "\nAll done ^_^"

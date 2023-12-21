#!/bin/bash

# Usage check
if [ "$#" -lt 4 ]; then
    echo "Usage: $0 <war-file-or-url> <groupId:artifactId:version> [<groupId:artifactId:version> ...]"
    exit 1
fi

# Name of this script
script_name=$(basename "$0")

# Clean up everything in the current directory except this script
find . -mindepth 1 ! -name "$script_name" -exec rm -rf {} \;

# Assign the first argument as the WAR file location or URL
WAR_FILE_OR_URL=$1
WAR_FILE_NAME="new.war"
# Remove the first argument so the rest can be treated as libraries
shift

# Function to handle the WAR file, whether it's a URL or a local file
handle_war_file() {
    

    if [[ $WAR_FILE_OR_URL == http* ]]; then
        echo "Downloading WAR file from URL: $WAR_FILE_OR_URL"
        curl -o $WAR_FILE_NAME $WAR_FILE_OR_URL
    else
        echo "Copying WAR file from local path: $WAR_FILE_OR_URL"
        cp "$WAR_FILE_OR_URL" $WAR_FILE_NAME
    fi

    if [ ! -f $WAR_FILE_NAME ]; then
        echo "Error: Failed to obtain WAR file."
        exit 1
    fi
}

handle_war_file

download_library_from_maven() {
    local group_id=$(echo $1 | cut -d: -f1)
    local artifact_id=$(echo $1 | cut -d: -f2)
    local version=$(echo $1 | cut -d: -f3)
    local group_id_path=$(echo $group_id | tr '.' '/')
    local maven_url="https://repo1.maven.org/maven2/${group_id_path}/${artifact_id}/${version}/${artifact_id}-${version}.jar"

    # Download the file using curl, and capture HTTP status code
    http_status=$(curl -o "${artifact_id}-${version}.jar" -w "%{http_code}" -s $maven_url)

    if [ "$http_status" -eq 200 ]; then
        echo "Downloaded ${artifact_id}-${version}.jar"
    else
        echo "Error: Library not found or download failed from Maven Central: $maven_url"
        rm -f "${artifact_id}-${version}.jar"  # Clean up partial or empty file
        exit 1
    fi
}

# Unwrap (unzip) cas.war
unzip $WAR_FILE_NAME

if [ ! -d "WEB-INF/lib" ]; then
    echo "Error: WEB-INF/lib directory not found in the WAR file."
    exit 1
fi

# Loop through each library specified in the arguments
for lib in "$@"; do
    download_library_from_maven $lib

    # Extract artifactId and version from the argument
    artifact_id=$(echo $lib | cut -d: -f2)
    version=$(echo $lib | cut -d: -f3)

    # find and remove the library under WEB-INF/lib
    find "WEB-INF/lib/" -type f -name "${artifact_id}-*.jar" -exec rm {} \;

    # Copy the new library to the correct location
    cp "${artifact_id}-${version}.jar" "WEB-INF/lib/"

    # Clean up downloaded library file
    rm "${artifact_id}-${version}.jar"
done

# Move script outside in order to not be included in war file
mv update-libs.sh ../update-libs.sh

# Recreate cas.war with new content
jar -cvf0m cas.war ./META-INF/MANIFEST.MF .

# Return script back in correct location
mv ../update-libs.sh update-libs.sh 

# Comment the following line if you want to keep the extracted files
rm -rf META-INF WEB-INF org

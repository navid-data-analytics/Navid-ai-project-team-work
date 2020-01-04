function make_release_tag() {
  git status > /dev/null 2>&1
  if [ $? -ne 0 ]; then
    echo "Current directory is not a git repository."
    return
  fi

  local git_tag=""
  if [ ! -z $1 ]; then
    git_tag=$(date "+v%Y%m%d-%H%M%SEET-$1")
  else
    git_tag=$(date "+v%Y%m%d-%H%M%SEET")
  fi

  git tag $git_tag
  echo "Created tag $git_tag"
  echo "Pushing to remote..."
  git push --tags
  echo "Done!"
}
make_release_tag $1

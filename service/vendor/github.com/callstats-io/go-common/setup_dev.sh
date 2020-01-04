echo 'Setting up git hooks...'
cd .git/hooks
ln -sf ../../git-hooks/pre-commit.sh ./pre-commit
ln -sf ../../git-hooks/prepare-commit-msg.sh ./prepare-commit-msg
ln -sf ../../git-hooks/pre-push.sh ./pre-push
echo 'Done'

echo 'Checking docker...'
docker ps >& /dev/null
if [ $? != 0 ]; then
  echo 'Could not detect a running docker machine.'
  echo 'You should setup docker before starting development.'
fi
echo 'Done'

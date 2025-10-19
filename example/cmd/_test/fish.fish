function _example_quote_suffix
  if not commandline -cp | xargs echo 2>/dev/null >/dev/null
    if echo (commandline -cp)'"' | xargs echo 2>/dev/null >/dev/null
      echo '"'
    else if echo (commandline -cp)"'" | xargs echo 2>/dev/null >/dev/null
      echo "'"
    end
  else 
    echo ""
  end
end

function _example_callback
  echo (commandline -cp)(_example_quote_suffix) | sed "s/ \$/ ''/" | xargs example _carapace fish
end

complete -e 'example'
complete -c 'example' -f -a '(_example_callback)' -r


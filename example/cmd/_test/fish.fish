function _example_quote_suffix
  if not commandline -cp | xargs echo 2>/dev/null >/dev/null
    if commandline -cp | sed 's/$/"/'| xargs echo 2>/dev/null >/dev/null
      echo '"'
    else if commandline -cp | sed "s/\$/'/"| xargs echo 2>/dev/null >/dev/null
      echo "'"
    end
  else 
    echo ""
  end
end

function _example_callback
  commandline -cp | sed "s/\$/"(_example_quote_suffix)"/" | sed "s/ \$/ ''/" | xargs example _carapace fish
end

complete -c example -f
complete -c 'example' -f -a '(_example_callback)' -r


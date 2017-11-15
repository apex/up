for src in dist/*; do
  dst=$(echo $src | sed 's/-pro//')
  mv $src $dst
done

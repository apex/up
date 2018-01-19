for src in dist/up_*; do
  dst=$(echo $src | sed 's/-pro//')
  mv $src $dst
done

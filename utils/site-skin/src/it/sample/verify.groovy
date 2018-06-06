File index = new File( basedir, "target/site/index.html" )
assert index.exists()

String html = index.getText()
assert html.contains( '<!DOCTYPE html PUBLIC' )

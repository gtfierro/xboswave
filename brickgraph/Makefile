generate:
	python make_graph.py

viz: generate
	python ontology-visualization/ontology_viz.py -o test.dot test.ttl
	dot -Tpdf -o test.pdf test.dot


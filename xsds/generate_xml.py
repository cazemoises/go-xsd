import xmlschema
import xml.etree.ElementTree as ET

def create_element_from_xsd(schema, element_name, namespace, depth=0):
    """Cria elementos XML recursivamente a partir de um XSD com namespace."""
    if depth > 10:  # Evita loops infinitos em XSDs complexos
        return None

    full_element_name = f"{{{namespace}}}{element_name}"
    if full_element_name not in schema.elements:
        raise ValueError(f"Elemento {element_name} não encontrado no XSD.")

    element = schema.elements[full_element_name]
    xml_element = ET.Element(element.name)

    if hasattr(element.type, 'content') and element.type.content:
        for child in element.type.content.iter_elements():
            child_element = create_element_from_xsd(schema, child.name, namespace, depth + 1)
            if child_element is not None:
                xml_element.append(child_element)

    return xml_element

# Caminho para o arquivo XSD
xsd_file = 'ACCC471.xsd'

# Define o namespace a partir do XSD
namespace = "http://www.cip-bancos.org.br/ARQ/ACCC471.xsd"

# Carrega o esquema XML
schema = xmlschema.XMLSchema(xsd_file)

# Gera um XML baseado no elemento raiz do XSD
root_element_name = 'ACCCDOC'  # Nome do seu elemento raiz
root = create_element_from_xsd(schema, root_element_name, namespace)

# Converte o XML gerado para uma string formatada
xml_string = ET.tostring(root, encoding='unicode', method='xml')

# Salva o XML gerado em um arquivo
with open('generated_example.xml', 'w') as f:
    f.write(xml_string)

print(f'XML gerado salvo em: generated_example.xml')

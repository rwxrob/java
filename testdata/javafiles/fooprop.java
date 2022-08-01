import java.util.Properties;
import java.util.Enumeration;

class Props {
    public static void main(String[] args) {
      Properties p = System.getProperties();
      String value = (String)p.get("foo");
      System.out.println(value);
    }
}

